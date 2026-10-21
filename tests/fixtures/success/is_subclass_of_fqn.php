<?php
interface LocalContract {}
class LocalBase implements LocalContract {}
class LocalChild extends LocalBase {}

echo (int) is_subclass_of(LocalBase::class, '\\LocalContract'), '|';
echo (int) is_subclass_of(LocalChild::class, '\\LocalBase'), '|';
echo (int) is_subclass_of(LocalBase::class, '\\LocalBase'), '|';
echo (int) is_subclass_of(ArrayIterator::class, '\\Traversable'), '|';
echo (int) is_a(new LocalChild, '\\LocalContract'), '|';
echo (int) is_a('\\LocalChild', '\\LocalBase', true);
