<?php
$a = 0; $o = new stdClass; $o->x = 1; $list = [];
false && ($a = 9); true || ($o->x = 9);
if (true) echo 'ok'; elseif ($list[] = 1 / 0) echo 'bad';
if (false) $list[] = $missing; else if (false) $o->x = 1 / 0;
echo $a, $o->x, ($list === []) + 0, "\n";
