/**
 * Determines whether the given string is a valid date. Returns a boolean.
 */
export function isValidDateString(str) {
    return !isNaN(Date.parse(str));
}

/**
 * Takes a valid date object and converts it to a string in the format of 'YYYY/MM/DD'.
 *
 * See: https://stackoverflow.com/questions/2013255/how-to-get-year-month-day-from-a-date-object
 */
export function convertDateToYYYYMMDD(date) {
    const month   = date.getUTCMonth() + 1; // months from 1-12
    const day     = date.getUTCDate();
    const year    = date.getUTCFullYear();
    // Using padded values, so that 2023/1/7 becomes 2023/01/07
    const pMonth        = month.toString().padStart(2,"0");
    const pDay          = day.toString().padStart(2,"0");
    return `${year}/${pMonth}/${pDay}`;
}
