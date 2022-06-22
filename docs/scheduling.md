### scheduling

*this scheduling architecture was written by someone who has never done a lick of scheduling outside of a science fair project, if you have a better implementation, feel free to step right up*

#### inputs

- time
    - `time start` - the starting time
    - `time divisions` - a list of durations specifying "slots of time"

- `students` - a list of students
- `judges` - a list of judges
- `events` - a list of events
- `rooms` - a list of occupiable rooms

- student requests - a list of "requests" by a student to join an event

#### outputs

- `assignment` - a student + the event they wish to attend
- `judgement` - a judge + the `assignments` they will judge throughout the divisions of time
- `housing` - a room + the `judgements` that will be happening in the room

#### the algorithm

the algorithm of the scheduler works bottom to top (from the most granular decisions to the larger ones)

1. parse the `student requests` into `assignments`.
1. assign requests to judges based on a few rules
    1. add events to judges that already have the same event type and are not full
    1. if there is no judge with the same event type, add it to a judge with no `assignments`
    1. if there are no judges with no `assignments`, keep it on hold
1. get the leftovers and assign them to judges that aren't full regardless of event type
    - if this is not possible, warn the user
1. attempt to spread the judges evenly across the rooms
