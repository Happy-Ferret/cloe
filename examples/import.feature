Feature: Import statement
  Scenario: Import a module
    Given a file named "main.tisp" with:
    """
    (import "mod")
    (seq ..(mod.map write [1 2 3 4 5]))
    """
    And a file named "mod.tisp" with:
    """
    (def (map func list)
        (match list
            [] []
            [first ..rest] (prepend (func first) (map func rest))))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    1
    2
    3
    4
    5
    """

  Scenario: Import a module in a directory
    Given a file named "main.tisp" with:
    """
    (import "modules/mod")
    (mod.Hello "world")
    """
    And a file named "modules/mod.tisp" with:
    """
    (def (Hello name) (write (merge "Hello, " name "!")))
    """
    When I successfully run `tisp main.tisp`
    Then the stdout should contain exactly:
    """
    Hello, world!
    """
