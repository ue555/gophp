<?php
class IncludedReadonly {
    public function __construct(public readonly int $value) {}
    public function change() { $this->value = 99; }
}
