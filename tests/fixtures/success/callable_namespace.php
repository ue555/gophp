<?php
namespace CallableExample;
function square(int $value): int { return $value * $value; }
$local = square(...);
$builtin = \strlen(...);
echo $local(4), '|', $builtin('xy'), '|';
echo (int)\is_callable('CallableExample\\square');
