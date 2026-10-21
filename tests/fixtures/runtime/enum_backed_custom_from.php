<?php
echo 'before';
enum BackedChoice: int {
    case A = 1;
    public static function from(int $value): self { return self::A; }
}
echo 'after';
