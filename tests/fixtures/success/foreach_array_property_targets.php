<?php
$out = []; $keys = []; foreach ([4 => 10, 's' => 20] as $keys[] => $out[]) {}
echo $keys[0], ':', $keys[1], ':', $out[0], ':', $out[1], "\n";
$o = new stdClass; foreach ([7] as $o->key => $o->value) {} echo $o->key, ':', $o->value, "\n";
