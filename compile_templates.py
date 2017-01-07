from local.template_conf import CONF, default_extensions

import os
import jinja2
import shutil


class TrackedLoader(jinja2.FileSystemLoader):
    """Template loader that tracks which files are actually loaded."""
    def __init__(self, *args, **kwargs):
        self.names_loaded = set()
        super().__init__(*args, **kwargs)

    def load(self, environment, name, *args, **kwargs):
        self.names_loaded.add(name)
        return super().load(environment, name, *args, **kwargs)


def conditional_write(path, content):
    """Overwrite the file at path with content, unless it is up-to-date."""
    try:
        with open(path, 'r') as outfile:
            if outfile.read() == content:
                return
    except FileNotFoundError:
        pass

    with open(path, 'w') as outfile:
        outfile.write(content)


default_extensions |= {
        ".html", ".css", ".js",
        ".eot", ".ttf", ".woff", ".woff2",
        ".png", ".jpg", ".jpeg", ".gif", ".bmp",
        ".mp3", ".mp4", ".mpeg", ".flac",
        ".webm", ".avi", ".wmv", ".mov", ".qt", ".flv",
        ".txt", ".tex", ".pdf",
}


def compile_templates():
    template_path = "local/templates"
    loader = TrackedLoader(template_path)
    env = jinja2.Environment(loader=loader, trim_blocks=True, lstrip_blocks=True)
    env.globals['global'] = CONF

    for name, local in CONF['pages'].items():
        output = env.get_template(name).render(current=name, local=local)

        path = os.path.join("static", name)
        directory = os.path.dirname(path)
        if not os.path.exists(directory):
            os.makedirs(directory)

        conditional_write(path, output)

    # Copy any files with default extensions not used in templates to static/, verbatim.
    for dirpath, _, filenames in os.walk(template_path):
        for filename in filenames:
            relsrc = os.path.join(os.path.relpath(dirpath, template_path), filename)
            if relsrc.startswith("./"):
                relsrc = relsrc[2:]
            _, ext = os.path.splitext(filename)
            if ext in default_extensions and relsrc not in loader.names_loaded:
                dstdir = os.path.join("static", os.path.relpath(dirpath, template_path))
                if not os.path.exists(dstdir):
                    os.makedirs(dstdir)
                src = os.path.join(dirpath, filename)
                dst = os.path.join(dstdir, filename)
                shutil.copy2(src, dst)


if __name__ == "__main__":
    compile_templates()
