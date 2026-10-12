<?php
$a = new stdClass; $b = new stdClass; $c = $a;
echo ($a == $b) + 0, ($a === $b) + 0, ($a === $c) + 0;
$a->x = 1; $b->x = '1';
echo ($a == $b) + 0, ($a !== $b) + 0, $c->x;
$c->x = 2; echo $a->x;
$a->self = $a; echo ($a === $c) + 0, ($a == $c) + 0, "\n";
