<?php
// ReflectionClass::getNamespaceName()/inNamespace() and
// ReflectionProperty::isPromoted(), used by Symfony TypeInfo and PropertyInfo.
namespace App\Models;

final class Book
{
    public string $isbn = '';

    public function __construct(public string $title, private readonly int $year = 2026)
    {
    }
}

$class = new \ReflectionClass(Book::class);
echo $class->getName(), '|', $class->getNamespaceName(), '|', $class->getShortName(), '|', $class->inNamespace() ? 'in' : 'global', "\n";
$global = new \ReflectionClass(\ArrayIterator::class);
echo $global->getNamespaceName() === '' ? 'empty' : 'set', '|', $global->inNamespace() ? 'in' : 'global', "\n";
foreach ($class->getProperties() as $property) {
    echo $property->getName(), '|', $property->isPromoted() ? 'promoted' : 'declared', '|', $property->isReadOnly() ? 'readonly' : 'mutable', "\n";
}
