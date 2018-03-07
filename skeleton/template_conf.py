from templating.loader import TrackedLoader
from templating.compiler import default_extensions


config = {
    'pages': {}
}

extensions_to_copy = default_extensions

loader = TrackedLoader("local/templates")
