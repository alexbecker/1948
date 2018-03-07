from local.template_conf import loader, extensions_to_copy, config
from templating.compiler import compile_templates


if __name__ == "__main__":
    compile_templates(loader, extensions_to_copy, config)
