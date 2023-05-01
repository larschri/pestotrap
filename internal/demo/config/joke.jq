.[] | . +
	{
		# Field `id` is required. It must be a string.
		id: .id | tostring,

		# Field `title` is used to display matches.
		# It is already present in the source data.

		# Field `description` is displayed below the title in the search results
		description: .category
	}
