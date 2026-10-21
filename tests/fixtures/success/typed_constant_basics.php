<?php
interface TypedContract { public const int VALUE = 1; }
trait TypedConstantTrait { public const string LABEL = 'trait'; }
class TypedConstantBase implements TypedContract {
    use TypedConstantTrait;
    public const int VALUE = 2;
    protected const int SECRET = 4;
    public static function secret() { return self::SECRET; }
}
class TypedConstantChild extends TypedConstantBase {
    public const int VALUE = 3;
}
echo TypedConstantChild::VALUE, '|', TypedConstantChild::LABEL, '|', TypedConstantChild::secret();
