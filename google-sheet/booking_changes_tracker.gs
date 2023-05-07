/**
 * Track booking changes and log affected rows data.
 * @param {Object} e - The event object from the onEdit trigger.
 */
function trackBookingChanges(e) {
    if (!e) {
        return;
    }

    const sheet = e.source.getSheetByName('db');
    const range = e.range;
    const affectedRows = getAffectedRows(sheet, range);

    if (!affectedRows || affectedRows.length === 0) {
        return;
    }

    const bookingsRows = processAffectedRows(affectedRows);
    Logger.log('Affected rows data: %s', JSON.stringify(bookingsRows, null, 2));
}

/**
 * Get the affected rows data from the sheet.
 * @param {Sheet} sheet - The sheet where changes occurred.
 * @param {Range} range - The range of edited cells.
 * @returns {Array} - The affected rows data.
 */
function getAffectedRows(sheet, range) {
    const startRow = range.getRow();
    const numRows = range.getNumRows();
    const totalCols = sheet.getLastColumn();

    return sheet.getRange(startRow, 1, numRows, totalCols).getValues();
}

/**
 * Process the affected rows data and create bookings rows.
 * @param {Array} affectedRows - The affected rows data.
 * @returns {Array} - The bookings rows data.
 */
function processAffectedRows(affectedRows) {
    const bookingsRows = [];

    affectedRows.forEach(row => {
        const date = row[0];

        if (!(date instanceof Date)) {
            return;
        }

        const rowData = row.slice(1);
        const bookings = extractBookings(rowData);

        bookingsRows.push({
            date: date, bookings: bookings,
        });
    });

    return bookingsRows;
}

/**
 * Extract bookings data from the row data.
 * @param {Array} rowData - The row data.
 * @returns {Array} - The extracted bookings data.
 */
function extractBookings(rowData) {
    const bookings = [];

    for (let i = 0; i < rowData.length; i += 2) {
        const clientName = rowData[i];
        const timeSlot = rowData[i + 1];

        if (clientName !== '' && timeSlot !== '') {
            bookings.push({
                clientName: clientName, timeSlot: timeSlot,
            });
        }
    }

    return bookings;
}
