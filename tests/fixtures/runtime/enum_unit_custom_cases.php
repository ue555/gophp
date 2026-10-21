<?php
echo 'before';
enum UnitChoice {
    case A;
    public static function cases(): array { return []; }
}
echo 'after';
