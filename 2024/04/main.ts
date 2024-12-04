import { fetchPuzzleInputLines } from '../utils.ts';
import { countCrossMASOccurrences, countGridOccurrences } from './lib.ts';

const lines = await fetchPuzzleInputLines(2024, 4);

const occurences = countGridOccurrences('XMAS', lines);

console.log('Part 1: There are', occurences, 'occurrences of XMAS');

const crossMASCount = countCrossMASOccurrences(lines);

console.log('Part 1: There are', crossMASCount, 'occurrences of X-MAS');
