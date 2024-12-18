import { getCharAtPoint, getGridSize, Point, walkGrid } from '../grids.ts';

const countOccurrencesInColumns = (rows: string[], target: string): number => {
  const getChar = getCharAtPoint(rows);
  const size = getGridSize(rows);

  let count = 0;
  for (let x = 0; x < size; x++) {
    let column = '';
    for (let y = 0; y < size; y++) {
      const char = getChar([x, y]);

      column += char;
    }

    count += countOccurrencesBothWays(target, column);
  }

  return count;
};

const countOccurrencesInRows = (rows: string[], target: string): number => {
  return sum(rows.map((row) => countOccurrencesBothWays(target, row)));
};

const countOccurrencesInDiagonals = (rows: string[], target: string) => {
  return getUpRightDiagonals(rows, target) +
    getDownRightDiagonals(rows, target);
};

const getUpRightDiagonals = (
  rows: string[],
  target: string,
): number => {
  const size = getGridSize(rows);

  const startingPoints: Point[] = [];
  for (let y = 0; y < size; y++) {
    startingPoints.push([0, y]);
  }

  for (let x = 1; x < size; x++) {
    startingPoints.push([x, size - 1]);
  }

  const counts = startingPoints.map((start) => {
    const path = walkGrid(rows, start, ([x, y]) => [x + 1, y - 1]);

    return countOccurrencesBothWays(target, path);
  });

  return sum(counts);
};

const getDownRightDiagonals = (
  rows: string[],
  target: string,
): number => {
  const size = getGridSize(rows);

  const startingPoints: Point[] = [];
  for (let y = 0; y < size; y++) {
    startingPoints.push([0, y]);
  }

  for (let x = 1; x < size; x++) {
    startingPoints.push([x, 0]);
  }

  const counts = startingPoints.map((start) => {
    const path = walkGrid(rows, start, ([x, y]) => [x + 1, y + 1]);

    return countOccurrencesBothWays(target, path);
  });

  return sum(counts);
};

const reverseString = (input: string): string => {
  return input.split('').toReversed().join('');
};

const countMatches = (target: string, text: string): number => {
  let index = text.indexOf(target, 0);

  let matches = 0;
  while (index !== -1) {
    matches++;
    index = text.indexOf(target, index + target.length);
  }

  return matches;
};

const countOccurrencesBothWays = (target: string, text: string) => {
  return countMatches(target, text) + countMatches(reverseString(target), text);
};

export const countXMASOccurrences = (
  rows: string[],
): number => {
  const target = 'XMAS';

  const rowCount = countOccurrencesInRows(rows, target);
  const columnCount = countOccurrencesInColumns(rows, target);
  const diagonalsCount = countOccurrencesInDiagonals(rows, target);

  return rowCount + columnCount + diagonalsCount;
};

const getCross = (rows: string[], centre: Point) => {
  const getChar = getCharAtPoint(rows);

  const centreChar = getChar(centre);

  const [x, y] = centre;
  return [
    `${getChar([x - 1, y - 1])}${centreChar}${getChar([x + 1, y + 1])}`,
    `${getChar([x - 1, y + 1])}${centreChar}${getChar([x + 1, y - 1])}`,
  ];
};

const sum = (numbers: number[]) => numbers.reduce((sum, n) => sum + n, 0);

const checkForCrossMAS = (rows: string[], coord: Point): boolean => {
  const cross = getCross(rows, coord);

  const counts = cross
    .map((diagonal) => countOccurrencesBothWays('MAS', diagonal));

  return sum(counts) === 2;
};

export const countCrossMASOccurrences = (rows: string[]): number => {
  const size = getGridSize(rows);

  let count = 0;
  for (let y = 1; y < size - 1; y++) {
    for (let x = 1; x < size - 1; x++) {
      if (rows[y][x] === 'A') {
        if (checkForCrossMAS(rows, [x, y])) count++;
      }
    }
  }

  return count;
};
