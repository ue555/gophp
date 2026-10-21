<?php
class LazyReadonly {
    public readonly ?int $value;
    public function fill() { $this->value ??= 4; }
}
$o = new LazyReadonly;
echo isset($o->value) ? 'set|' : 'unset|';
echo empty($o->value) ? 'empty|' : 'filled|';
echo $o->value ?? 'fallback', '|';
$o->fill();
$o->fill();
echo $o->value, "\n";
