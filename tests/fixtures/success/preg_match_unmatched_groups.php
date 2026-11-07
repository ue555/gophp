<?php
// Unmatched capture groups: "" when a later group matched, dropped when
// trailing, null with PREG_UNMATCHED_AS_NULL, [value, -1] with
// PREG_OFFSET_CAPTURE; named groups come before their index.
// Symfony Finder's NumberComparator reads $matches[1] ?: '==' for '> 0' and '0'.
foreach (['> 0', '0', '== 1k', ' 5 '] as $test) {
    preg_match('#^\s*(==|!=|[<>]=?)?\s*([0-9\.]+)\s*([kmg]i?)?\s*$#i', $test, $matches);
    echo json_encode($matches), '|', $matches[1] ?: '==', "\n";
}
preg_match('/(a)?(b)?(c)/', 'c', $matches);
echo json_encode($matches), "\n";
preg_match('/(a)?(b)?(c)/', 'c', $matches, PREG_UNMATCHED_AS_NULL);
echo json_encode($matches), "\n";
preg_match('/(a)?(b)?(c)/', 'xc', $matches, PREG_OFFSET_CAPTURE);
echo json_encode($matches), "\n";
preg_match('/(a)?(b)?(c)/', 'xc', $matches, PREG_OFFSET_CAPTURE | PREG_UNMATCHED_AS_NULL);
echo json_encode($matches), "\n";
preg_match('/(?<x>a)?(b)(?<y>c)?/', 'b', $matches);
echo json_encode($matches), "\n";
preg_match('/(?<x>a)?(b)(?<y>c)?/', 'b', $matches, PREG_UNMATCHED_AS_NULL);
echo json_encode($matches), "\n";
preg_match('/(?<word>\w+)/', 'hi', $matches);
echo json_encode($matches), "\n";
