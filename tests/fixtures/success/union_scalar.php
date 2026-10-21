<?php
function scalarPick(int|string $value): int|string { return $value; }
function numberPick(float|int $value): int|float { return $value; }
echo gettype(scalarPick(7)), '|', gettype(scalarPick('7')), '|';
echo gettype(numberPick('7')), '|', gettype(numberPick('7.5')), '|';
echo gettype(numberPick(true)), '|', gettype(numberPick(2.5)), "\n";
