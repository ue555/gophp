<?php
trait CustomTryFrom {
    public static function tryFrom(int $value): ?self { return null; }
}
echo 'before';
enum BackedTraitChoice: int {
    use CustomTryFrom;
    case A = 1;
}
echo 'after';
