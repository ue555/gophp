<?php
readonly class ReadonlyBase {
    public function __construct(public int $value) {}
}
readonly class ReadonlyChild extends ReadonlyBase {}
$child = new ReadonlyChild(3);
echo $child->value, '|';
try { $child->value = 4; } catch (Error $e) { echo 'readonly'; }
