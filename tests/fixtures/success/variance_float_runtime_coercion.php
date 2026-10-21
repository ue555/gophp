<?php
declare(strict_types=1);
function floatArgument(float $number): float { return $number; }
function floatReturn(): float { return 1; }
echo gettype(floatArgument(1)), '|', gettype(floatReturn()), "\n";
