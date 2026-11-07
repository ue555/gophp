<?php
// Every Reflection* object is a Reflector (which extends Stringable), so it
// satisfies Reflector-typed parameters such as Symfony TypeInfo's
// TypeContextFactory::createFromReflection(Reflector $reflection).
final class Book
{
    public string $title = '';

    public function describe(int $width = 80): string
    {
        return $this->title;
    }
}

function describe(Reflector $reflection): string
{
    return get_class($reflection);
}

$class = new ReflectionClass(Book::class);
$method = $class->getMethod('describe');
$parameter = $method->getParameters()[0];
$property = new ReflectionProperty(Book::class, 'title');
$function = new ReflectionFunction('describe');
$type = $method->getReturnType();
foreach ([$class, $method, $parameter, $property, $function] as $reflection) {
    echo describe($reflection), '|', $reflection instanceof Reflector ? 'Reflector' : 'no', '|', $reflection instanceof Stringable ? 'Stringable' : 'no', "\n";
}
// ReflectionType objects are Stringable but not Reflector.
echo get_class($type), '|', $type instanceof Reflector ? 'Reflector' : 'no', '|', $type instanceof Stringable ? 'Stringable' : 'no', "\n";
$implemented = class_implements(ReflectionMethod::class);
sort($implemented); // builtin registration order differs between engines; the set is what matters
echo (new ReflectionClass(Reflector::class))->isInterface() ? 'interface' : 'class', '|', implode(',', $implemented), "\n";
