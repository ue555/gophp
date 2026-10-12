<?php
$o = new stdClass(); $o->name = 'PHP'; $o->Name = 'GO';
$o->data = ['x' => 4]; $o->data['x'] = 5; $o->list[] = 6;
$o->child = new STDCLASS; $o->child->v = 7;
echo $o->name, $o->Name, $o->data['x'], $o->list[0], $o->child->v, "\n";
echo (($o->v = 8) === 8) + 0, "\n";
