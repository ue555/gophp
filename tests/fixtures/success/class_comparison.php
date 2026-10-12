<?php
class One { private $x = 1; public $y = [2]; }
class Two { public $y = [2]; private $x = 1; }
$a = new One; $b = new One; $c = new Two;
echo ($a == $b) + 0, ($a === $b) + 0, ($a <> $c) + 0;
echo ($a < $c) + 0, ($a > $c) + 0, ($a <= $c) + 0, ($a >= $c) + 0;
echo ':', $a <=> $c, ':', $c <=> $a, "\n";
