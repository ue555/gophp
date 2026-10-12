<?php
$o = new stdClass; $o->n = 1;
$a = ['object' => $o, 'array' => [2]]; $b = $a;
$b['object']->n = 9; $b['array'][0] = 3;
echo $a['object']->n, $o->n, $a['array'][0], $b['array'][0];
$o->data = [4]; $c = $o->data; $c[0] = 5;
echo $o->data[0], $c[0], "\n";
