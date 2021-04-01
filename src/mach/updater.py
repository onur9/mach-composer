import os
import re
from abc import ABC
from dataclasses import dataclass
from typing import List, Optional, Tuple, Union

import click
from mach import cache, exceptions, git
from mach.types import ComponentConfig, MachConfig

NAME_RE = re.compile(r".* name: [\"']?(.*)[\"']?")
VERSION_RE = re.compile(r"(\s*version: )([\"']?.*[\"']?)")

ROOT_BLOCK_START = re.compile(r"^\w+")


Updates = List[Tuple[ComponentConfig, str]]


@dataclass
class UpdaterInput:
    data: Union[MachConfig, List[ComponentConfig]]
    file: str

    @property
    def is_mach_config(self):
        return isinstance(self.data, MachConfig)


def update_config_component(  # noqa: C901
    updater_input: UpdaterInput,
    component_name: str,
    new_version: str,
):
    config = updater_input.data

    component = config.get_component(component_name)
    if not component:
        raise exceptions.MachError(f"Could not find component {component_name}")

    if component.version == new_version:
        click.echo(f"Component {component_name} is already on version {new_version}.")

    click.echo(f"Updating {component_name} to version {new_version}...")

    FileUpdater.apply_updates(updater_input.file, [(component, new_version)])


def update_config_components(  # noqa: C901
    updater_input: UpdaterInput,
    *,
    check_only=False,
    verbose=False,
):
    """
    Update a given MACH configuration file.

    :param config: The MACH configuration to update components for
    :param check_only: Only check for updates; don't update the file
    :param verbose: Enable verbose output
    """
    intro_msg = f"Checking updates for components in {updater_input.file}"
    print(intro_msg)
    print("-" * len(intro_msg))

    updates: Updates = _fetch_changes(updater_input)

    updater_cls = FileUpdater if updater_input.is_mach_config else ComponentsFileUpdater

    if not check_only:
        updater_cls.apply_updates(updater_input.file, updates)


def _fetch_changes(updater_input: UpdaterInput) -> Updates:
    cache_dir = cache.cache_dir_for(updater_input.file)

    if updater_input.is_mach_config:
        components = updater_input.data.components
    else:
        components = updater_input.data

    updates: Updates = []
    for component in components:
        click.echo(f"Updates for {component.name}...")

        match = re.match(r"^git::(.*)", component.source)
        if not match:
            click.echo(
                f"  Cannot check {component.name} component since it doesn't have a Git source defined"  # noqa
            )
            continue

        component_dir = os.path.join(cache_dir, component.name)
        repo = match.group(1)
        match = re.match(r"(.*\/.*)(\/\/.*)$", repo)
        if match:
            repo = match.group(1)
        git.ensure_local(repo, component_dir)

        commits = git.history(component_dir, component.version, branch=component.branch)
        if not commits:
            click.echo("  No updates\n")
            continue

        for commit in commits:
            print(f"  {commit.id}: {commit.msg}")

        click.echo("")

        updates.append((component, commits[0].id))

    return updates


class BaseFileUpdater(ABC):
    """Updater which update component version in-place.

    We'll use a very basic search-and-replace based on regular expressions
    instead of the yaml parser to not mess with any formatting.
    """

    @classmethod
    def apply_updates(cls, file: str, updates: Updates):
        """Apply given updates to the file."""
        instance = cls()
        instance.apply(file, updates)

    def apply(self, file: str, updates: Updates):
        click.echo("Writing updated to file...")
        self.current_component: Optional[ComponentConfig] = None
        self.updates = {component.name: version for component, version in updates}
        self.component_map = {component.name: component for component, _ in updates}

        with open(file) as f:
            lines = f.readlines()

        newlines = [self.process_line(line) for line in lines]

        with open(file, mode="w") as f:
            for line in newlines:
                f.write(line)

    def process_line(self, line: str) -> str:
        raise NotImplementedError()

    def process_component_line(self, line: str):
        name_match = NAME_RE.match(line)
        if name_match:
            component_name = name_match.group(1)

            print(f"Processing {component_name}")

            try:
                self.current_component = self.component_map[component_name]
            except KeyError:
                self.current_component = None

            return line

        if not self.current_component:
            return line

        match = VERSION_RE.match(line)
        if not match:
            return line

        assert self.current_component.version in match.group(2)

        try:
            new_version = self.updates[self.current_component.name]
        except KeyError:
            return line

        if new_version.isdigit():
            new_version = f'"{new_version}"'

        return VERSION_RE.sub(rf"\g<1>{new_version}", line)


class FileUpdater(BaseFileUpdater):
    def __init__(self):
        self.in_components = False

    def process_line(self, line: str) -> str:
        if line.startswith("components:"):
            self.in_components = True
        elif ROOT_BLOCK_START.match(line):
            self.in_components = False
        elif self.in_components:
            return self.process_component_line(line)

        return line


class ComponentsFileUpdater(BaseFileUpdater):
    def process_line(self, line: str) -> str:
        return self.process_component_line(line)
