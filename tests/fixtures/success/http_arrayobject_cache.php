<?php
$collection = new ArrayObject(['name' => 'Book']);
$copy = unserialize(serialize($collection));
echo count($copy), '|', $copy['name'];
