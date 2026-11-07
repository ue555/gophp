<?php
// ReflectionParameter::getDeclaringFunction() returns the ReflectionMethod or
// ReflectionFunction that declares the parameter; getPosition() is its index.
// Symfony TypeInfo resolves parameter types through the declaring function.
final class Book
{
    public function __construct(public string $title, int $year = 2026)
    {
    }

    public static function make(string $title, int ...$years): self
    {
        return new self($title);
    }
}

function describe(Book $book, ?string $prefix = null): string
{
    return ($prefix ?? '') . $book->title;
}

foreach ([new ReflectionMethod(Book::class, '__construct'), new ReflectionMethod(Book::class, 'make'), new ReflectionFunction('describe')] as $function) {
    foreach ($function->getParameters() as $parameter) {
        $declaring = $parameter->getDeclaringFunction();
        echo get_class($declaring), '|', $declaring->getName(), '|', $parameter->getPosition(), '|', $parameter->getName(), '|', $declaring instanceof ReflectionMethod ? $declaring->getDeclaringClass()->getName() : '-', '|', $declaring === $function ? 'same' : 'copy', "\n";
    }
}
