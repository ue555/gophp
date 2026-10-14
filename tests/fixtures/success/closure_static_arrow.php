<?php
$double = static function (int $x): int { return $x * 2; };
$next = static fn (int $n): int => $n + 1;
echo $double(3), '|', $next(4), "\n";
