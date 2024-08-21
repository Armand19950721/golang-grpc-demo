exports.replaceColumn = (obj, old_col, new_col) => {
    if (!!obj[old_col]) {
        obj[new_col] = obj[old_col];
    }
    delete obj[old_col]
}