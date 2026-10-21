<?php
$o = (object) ['name' => 'box', 2 => 7];
$n = (object) null;
echo $o->name, '|', $o->{'2'}, '|', isset($n->scalar) ? 'set' : 'clear', '|', isset($_ENV['UNSET']) ? 'set' : 'clear';
