from os.path import splitext
from pathlib import Path
from typing import Dict, List

import click
import yaml
from mach import exceptions
from mach.types import ComponentConfig, MachConfig, SiteAzureSettings
from mach.validate import validate_config
from marshmallow.exceptions import ValidationError


def parse_configs(files: List[str], output_path: str = None) -> List[MachConfig]:
    """Parse and validate configurations."""
    valid_configs = []
    for file in files:
        config = parse_config_from_file(file)
        config.file = file
        click.echo(f"Parsed {file} into config")

        validate_config(config)

        config = resolve_general_config(config)
        config = resolve_components(config)

        if output_path:
            full_output_path = Path(f"{output_path}/{splitext(file)[0]}")
            full_output_path.mkdir(exist_ok=True, parents=True)
            config.output_path = str(full_output_path)

        valid_configs.append(config)
    return valid_configs


def parse_config_from_file(file: str) -> MachConfig:
    """Parse file into MachConfig object."""
    click.echo(f"Parsing {file}...")
    with open(file, "r") as fh:
        dictionary_config = yaml.full_load(fh)

    try:
        config = MachConfig.schema(infer_missing=True).load(dictionary_config)  # type: ignore
    except KeyError as e:
        # Most probably a missing value in the configuration.
        # dataclasses_json doesn't really give a proper Exception for this.
        # TODO: See if we can improve this / make it more robust. Either by improving
        # dataclassess_json (with a PR) or by extending it (if possible)
        raise exceptions.ParseError(f"Required attribute {e} missing") from e
    except ValidationError as e:
        # TODO: We don't have any path here, so not the best of error messages
        raise exceptions.ParseError(
            "Configuration file could not be validated", details=e.normalized_messages()
        ) from e
    return config


def resolve_general_config(config: MachConfig) -> MachConfig:
    """If no general config is specified, use global config settings."""
    if not config.general_config.azure:
        return config

    for site in config.sites:
        if site.azure:
            site.merge(config.general_config.azure)
        else:
            site.azure = SiteAzureSettings.from_config(config.general_config.azure)

    return config


def resolve_components(config: MachConfig) -> MachConfig:
    """If no component info is specified, use global component settings."""
    component_info: Dict[str, ComponentConfig] = {
        component.name: component for component in config.components
    }
    for site in config.sites:
        if site.components:
            for component in site.components:
                if not component.version:
                    component.version = component_info[component.name].version
                if not component.source:
                    component.source = component_info[component.name].source
                if not component.short_name:
                    component.short_name = component_info[component.name].short_name
    return config