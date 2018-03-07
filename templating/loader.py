import jinja2


class TrackedLoader(jinja2.FileSystemLoader):
    """Template loader that tracks which files are actually loaded."""
    def __init__(self, *args, **kwargs):
        self.names_loaded = set()
        super().__init__(*args, **kwargs)

    def get_source(self, environment, template):
        self.names_loaded.add(template)
        return super().get_source(environment, template)
