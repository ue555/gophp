<?php
class TypedBasics {
    public int $count;
    public ?string $label;
    public static int $shared;
}
$a = new TypedBasics;
$b = new TypedBasics;
echo isset($a->count) ? 'set' : 'unset';
echo '|', empty($a->count) ? 'empty' : 'full';
echo '|', $a->count ?? 'fallback';
$a->count ??= '12';
$a->count++;
$b->count = 9;
echo '|', $a->count, ':', $b->count;
echo '|', isset($a->label) ? 'set' : 'unset';
$a->label = null;
echo '|', isset($a->label) ? 'set' : 'unset';
TypedBasics::$shared = 4;
echo '|', TypedBasics::$shared;
