<?php
interface TypedA {}
interface TypedB {}
class TypedBoth implements TypedA, TypedB {}
class TypedUnion {
    public int|string $simple;
    public TypedA&TypedB $both;
    public (TypedA&TypedB)|null $optional;
}
$value = new TypedUnion;
$value->simple = 7;
$value->both = new TypedBoth;
$value->optional = $value->both;
echo gettype($value->simple), '|', $value->both instanceof TypedA, '|', $value->optional instanceof TypedB;
$value->simple = 'text';
$value->optional = null;
echo '|', $value->simple, '|', isset($value->optional) ? 'set' : 'unset';
