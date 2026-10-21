<?php declare(strict_types=1);
function unionStrict(int|float $value): int|float { return $value; }
echo gettype(unionStrict(3)), '|', gettype(unionStrict(4.5)), '|';
try { unionStrict('3'); } catch (TypeError $e) { echo 'TypeError'; }
echo "\n";
