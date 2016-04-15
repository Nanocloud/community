/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2016 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package errors

import "net/http"

var (
	InternalError = &apiError{
		0x000001,
		http.StatusInternalServerError,
		"An unexpected error occured.",
	}

	InvalidError = &apiError{
		0x000002,
		http.StatusInternalServerError,
		"An invalid error has been returned by the system.",
	}

	Unauthorized = &apiError{
		0x000003,
		http.StatusUnauthorized,
		"You are not authorized to perform this action.",
	}

	AdminLevelRequired = &apiError{
		0x000004,
		http.StatusUnauthorized,
		"Only administrators are allowed to perform this action.",
	}

	UserNotFound = &apiError{
		0x000005,
		http.StatusNotFound,
		"This user doesn't exist.",
	}

	InvalidRequest = &apiError{
		0x000006,
		http.StatusBadRequest,
		"The request is not valid.",
	}

	WindowsNotOnline = &apiError{
		0x000007,
		http.StatusServiceUnavailable,
		"Windows is not available.",
	}

)
