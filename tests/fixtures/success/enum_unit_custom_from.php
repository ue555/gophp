<?php
trait ChoiceLookup {
    public static function tryFrom(string $value): ?self {
        return $value === 'a' ? self::A : null;
    }
}
enum Choice {
    use ChoiceLookup;
    case A;
    public static function from(string $value): self { return self::A; }
}
echo Choice::from('x')->name, '|', Choice::tryFrom('a')->name, '|';
echo Choice::tryFrom('other') === null ? 'null' : 'bad';
