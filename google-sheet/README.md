# Google Sheet Booking Changes Tracker

This script is designed to track booking changes in a Google Sheet and output the affected rows with booking data.

## Google Sheet Format

The Google Sheet that the script parses should have the following format:

| Date              | Client Name 1 | Time Slot 1 | Client Name 2 | Time Slot 2 | ... |
|-------------------|---------------|-------------|---------------|-------------|-----|
| yyyy-mm-dd        | John Doe      | 10:00       | Jane Smith    | 12:00       | ... |

- The first column should contain the date (in the format yyyy-mm-dd).
- Subsequent columns should contain pairs of client names and time slots.
- Each client name should be followed by the corresponding time slot in the next column.
- You can add as many client name and time slot pairs as needed.

## Output Format

The script will output an array of objects, where each object represents an affected row in the Google Sheet. The objects will have the following format:

```json
{
  "date": "yyyy-mm-dd+00:00",
  "bookings": [
    {
      "client_name": "John Doe",
      "time": "10:00"
    },
    {
      "client_name": "Jane Smith",
      "time": "12:00"
    },
    ...
  ]
}
```

- The `date` property represents the date of the affected row (in the format yyyy-mm-ddT00:00:00.000Z).
- The `bookings` property is an array of objects that contain the client names and time slots of the affected bookings.
- Each object inside the `bookings` array has a `client_name` property and a `time` property.