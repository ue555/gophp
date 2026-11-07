<?php
// ReflectionProperty::$class names the declaring class, also for inherited and
// static properties; Symfony Serializer's AttributeLoader compares it with the
// reflected class name to skip inherited properties.
class Base
{
    public string $inherited = 'base';
    protected static int $count = 0;
}

final class Book extends Base
{
    public string $title = '';
    private ?int $year = null;

    public function describe(): string
    {
        return $this->title;
    }
}

$property = new ReflectionProperty(Book::class, 'title');
echo $property->class, '|', $property->name, "\n";
foreach ((new ReflectionClass(Book::class))->getProperties() as $reflected) {
    echo $reflected->class, '::$', $reflected->name, '|', $reflected->getDeclaringClass()->getName(), '|', $reflected->class === Book::class ? 'own' : 'inherited', "\n";
}
$method = new ReflectionMethod(Book::class, 'describe');
echo $method->class, '|', $method->name, "\n";
