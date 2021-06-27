def parse_val(value: str, fn: type = None):
    if value == '(none)':
        return None
    if fn:
        return fn(value)
    return value.strip()
