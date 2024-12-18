import { fetchPuzzleInputLines } from '../utils.ts';
import { countVisited, printVisited } from './lib.ts';

const grid = await fetchPuzzleInputLines(2024, 6);

console.log('Part 1: Number of visited spaces is', countVisited(grid));

printVisited(grid);
