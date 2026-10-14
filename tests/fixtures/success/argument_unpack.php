<?php
function combine($a, $b, $c) { return "$a$b$c"; }
function countValues(...$values) { return count($values); }
$values = [2, 3];
echo combine(1, ...$values), '|', countValues(...[4, 5]), "\n";
