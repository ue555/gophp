<?php
$a = new stdClass; $a->x = 1; $a->y = 2;
$b = new stdClass; $b->y = '2'; $b->x = 1;
echo ($a == $b) + 0, ($a === $b) + 0, ':', $a <=> $b, "\n";
$b->x = 3; echo ($a < $b) + 0, ($b > $a) + 0;
echo ($a > []) + 0, ($a > 'z') + 0, ($a <> false) + 0, "\n";
