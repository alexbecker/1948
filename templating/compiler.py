import os
import shutil

import jinja2

STATIC_PATH = "local/static"


def conditional_write(path, content):
    """
    Overwrite the file at path with content, unless it is up-to-date.

    Useful to prevent make from needlessly rebuilding targets.
    """
    try:
        with open(path, 'r') as outfile:
            if outfile.read() == content:
                return
    except FileNotFoundError:
        pass

    with open(path, 'w') as outfile:
        outfile.write(content)


default_extensions = {
        ".html", ".css", ".js",
        ".eot", ".ttf", ".woff", ".woff2",
        ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".ico",
        ".mp3", ".mp4", ".mpeg", ".flac",
        ".webm", ".avi", ".wmv", ".mov", ".qt", ".flv",
        ".txt", ".tex", ".pdf",
}


def compile_templates(loader, extensions_to_copy, config):
    env = jinja2.Environment(loader=loader, trim_blocks=True, lstrip_blocks=True)
    env.globals['global'] = config

    for name, local in config['pages'].items():
        output = env.get_template(name).render(current=name, local=local)

        path = os.path.join(STATIC_PATH, name)
        directory = os.path.dirname(path)
        if not os.path.exists(directory):
            os.makedirs(directory)

        conditional_write(path, output)

    # Copy any files with default extensions not used in templates to static/, verbatim.
    for template_path in loader.searchpath:
        for dirpath, _, filenames in os.walk(template_path):
            for filename in filenames:
                relsrc = os.path.join(os.path.relpath(dirpath, template_path), filename)
                if relsrc.startswith("./"):
                    relsrc = relsrc[2:]
                _, ext = os.path.splitext(filename)
                if ext in extensions_to_copy and relsrc not in loader.names_loaded:
                    dstdir = os.path.join(STATIC_PATH, os.path.relpath(dirpath, template_path))
                    if not os.path.exists(dstdir):
                        os.makedirs(dstdir)
                    src = os.path.join(dirpath, filename)
                    dst = os.path.join(dstdir, filename)
                    shutil.copy2(src, dst)
