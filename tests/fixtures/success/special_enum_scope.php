<?php
interface CaseContract { public function same(self $value): self; }
enum SpecialCase implements CaseContract {
    case One;
    public function same(CaseContract $value): self { return self::One; }
}
echo SpecialCase::One->same(SpecialCase::One) === SpecialCase::One ? 'enum' : 'wrong';
