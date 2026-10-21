<?php
interface First {}
interface Second extends First {}
class Item implements Second {}
$items = class_implements(Item::class);
echo $items['Second'], '|', $items['First'], '|';
echo strspn('abc123', 'abc'), '|';
echo unpack('C*', 'AZ')[2], '|';
echo str_repeat('x', 3), '|';
echo basename('/a/b/file.php', '.php');
