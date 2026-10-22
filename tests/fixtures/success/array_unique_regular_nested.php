<?php
$a = ['first' => ['src' => 'app.css', 'chunk' => ['v' => 1]], 'duplicate' => ['chunk' => ['v' => 1], 'src' => 'app.css'], 'different' => ['src' => 'app.js', 'chunk' => ['v' => 2]]];
echo json_encode(array_keys(array_unique($a, SORT_REGULAR))), "\n";
echo count($a), "\n";
$x = (object)['value' => 1]; $y = (object)['value' => 1];
echo json_encode(array_keys(array_unique([$x, $y, (object)['value' => 2]], SORT_REGULAR))), "\n";
echo json_encode(array_unique(['a', 'a', 'b'])), "\n";
try { array_unique([], flags: []); } catch (Throwable $e) { echo get_class($e), ': ', $e->getMessage(), "\n"; }
