def verticalSplit(data, scema):
    try:
        data = data.rstrip() # Remove the extra newline (if it exists)
        if "\n" in data:
            new_entries = []
            entries = data.split("\n")
            for entry in entries:
                new_entries.append(_makeEntry(entry, scema))
            return new_entries
        else:
            return _makeEntry(data, scema)
    except Exception:
        pass
        return None

        
def _makeEntry(data, scema):
    values = data.split("|")
    entry = {key: None for key in scema}
    for i in range(0, len(scema)):
        if scema[i].startswith("_"):
            del entry[scema[i]]
            continue
        entry[scema[i]] = values[i]

    return entry