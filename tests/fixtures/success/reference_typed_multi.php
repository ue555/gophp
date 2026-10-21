<?php
class TypedMulti {
    public ?int $number = null;
    public ?string $text = null;
}
$box = new TypedMulti;
$alias =& $box->number;
$box->text =& $alias;
try { $alias = '42'; } catch (TypeError $e) { echo 'conflict|'; }
echo isset($box->number) ? 'set' : 'null', '|';
unset($box->text);
$alias = '42';
echo $box->number, '|', gettype($alias);
