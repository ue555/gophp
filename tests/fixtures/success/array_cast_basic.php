<?php
$text = (array) 'one';
$empty = (array) null;
$object = new stdClass;
$object->name = 'box';
$fields = (array) $object;
echo $text[0], '|', count($empty), '|', $fields['name'];
